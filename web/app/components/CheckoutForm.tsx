import {
  PaymentElement,
  useElements,
  useStripe,
} from "@stripe/react-stripe-js";
import type { StripePaymentElementOptions } from "@stripe/stripe-js";
import { type FormEvent, useEffect, useState } from "react";
import { Disclaimer } from "./Disclaimer";
import { PaymentIntentError } from "./PiError";
import { useGenerateHoroscope } from "~/hooks/generateHoroscope";
import { Spinner } from "./Spinner";
import { AnimatePresence, easeInOut, motion } from "motion/react";

import axios from "axios";
import { Frown, Meh, Smile, SquareArrowOutUpRight, Copy, Check } from 'lucide-react';
import { useAddRating, type AddRatingRequest } from "~/hooks/addRating";
import { useShareHoroscope } from "~/hooks/shareHoroscope";

const FEEDBACK_DELAY = 2000 // TODO set back to 8



const getButtonColors = ({
  isLoading,
  isDisabled,
}: {
  isLoading: boolean;
  isDisabled: boolean;
}): GetButtonColorsResponse => {
  if (isLoading || isDisabled) {
    return {
      cursor: "cursor-not-allowed",
      background: "bg-pink-300",
      backgroundHover: "hover:bg-pink-300",
    };
  }

  return {
    cursor: "cursor-pointer",
    background: "bg-pink-500",
    backgroundHover: "hover:bg-pink-400",
  };
};

type CheckoutFormProps = {
  clientSecret: string;
};

const paymentElementOptions: StripePaymentElementOptions = {
  layout: "accordion",
};


type IconProps = {
  color: string
  icon: React.ReactNode
  callback: () => void
}
const Icon = ({ color, icon, callback }: IconProps) => {
  return (<div onClick={callback} className={`flex flex-col hover:cursor-pointer ${color} transition duration-150  scale-125 hover:scale-150`}>
    {icon}
  </div>
  )
}

// TODO click to send req




const Feedback = ({ horoscopeId, paymentIntentId }) => {
  const addRating = useAddRating()


  const handleRating = (addRatingRequest: AddRatingRequest) => {
    addRating.mutate(addRatingRequest)
  }


  return (
    <div className="flex max-w-sm justify-center items-center min-h-[140px]">
      {!addRating.isIdle ? <p className="text-center font-bold">Thanks for your feedback!</p> : (

        <div className="flex flex-col justify-center space-y-4">
          <p className="text-center font-bold ">How do you feel about your horoscope?</p>
          <div className="flex flex-row justify-around p-5 overflow-visible ">
            <Icon callback={() => handleRating({
              horoscopeId,
              paymentIntentId,
              rating: "negative"
            })} color="text-red-500" icon={<Frown />}></Icon>
            <Icon callback={() => handleRating({
              horoscopeId,
              paymentIntentId,
              rating: "neutral"
            })} color="text-yellow-500" icon={<Meh />}></Icon>
            <Icon callback={() => handleRating({
              horoscopeId,
              paymentIntentId,
              rating: "positive"
            })} color="text-emerald-500" icon={<Smile />}></Icon>
          </div>
        </div>

      )

      }
    </div>
  )
}



export function CheckoutForm({ clientSecret }: CheckoutFormProps) {
  const stripe = useStripe();
  const elements = useElements();
  const generateHoroscope = useGenerateHoroscope();
  const shareHoroscope = useShareHoroscope()
  const [errorMessage, setErrorMessage] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [isDisclaimerChecked, setIsDisclaimerChecked] = useState(false);
  const [feedbackVisible, setFeedbackVisible] = useState(false)
  const [paymentIntentId, setPaymentIntentId] = useState<null | string>(null)
  const [copied, setCopied] = useState(false)


  const shareableLink = `${window.location.origin}/${generateHoroscope?.data?.data?.external_id}`


  const isButtonDisabled =
    isLoading || !stripe || !elements || !isDisclaimerChecked;
  const button = getButtonColors({
    isLoading,
    isDisabled: !isDisclaimerChecked,
  });

  const handleCopy = async () => {
    await navigator.clipboard.writeText(shareableLink)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();

    // Check if disclaimer checkbox is checked
    if (!isDisclaimerChecked) {
      setErrorMessage(
        "Please accept the Terms of Service and Privacy Policy to continue."
      );
      return;
    }

    if (!stripe || !elements) {
      return;
    }

    setErrorMessage(""); // clear it
    setIsLoading(true);

    const { error: submitError } = await elements.submit();
    if (submitError) {
      setErrorMessage(submitError.message ?? "An unexpected error occurred");
      setIsLoading(false);
      return;
    }

    const { error, paymentIntent } = await stripe.confirmPayment({
      clientSecret,
      elements,
      confirmParams: {
        return_url: window.location.href,
      },
      redirect: "if_required",
    });

    if (error?.message) {
      if (error?.type === "card_error" || error?.type === "validation_error") {
        setErrorMessage(error.message);
        setIsLoading(false);
      } else {
        console.log(error);
        setErrorMessage("An unexpected error ocurred");
        setIsLoading(false);
      }
    }

    if (paymentIntent?.status === "succeeded") {
      localStorage.setItem("paymentIntentId", paymentIntent.id);
      // Generate the horoscope
      generateHoroscope.mutate(
        {
          paymentIntentId: paymentIntent.id,
        },
        {
          onSuccess: () => {
            setPaymentIntentId(paymentIntent.id)
            setIsLoading(false);
          },
          onError: (e) => {
            let message = "An unexpected error occurred";
            if (axios.isAxiosError(e) && e.response?.data?.message) {
              message = e.response?.data?.message;
            }
            setErrorMessage(message);
            setIsLoading(false);
          },
        }
      );
    } else {
      setIsLoading(false);
    }
  };

  const horoscope = generateHoroscope?.data?.data?.horoscope;


  useEffect(() => {
    if (horoscope) {
      const timer = setTimeout(() => {
        setFeedbackVisible(true)
      }, FEEDBACK_DELAY)

      return () => clearTimeout(timer)
    }
  }, [horoscope])

  return (
    <div className="relative overflow-visible">
      <AnimatePresence mode="wait">
        {horoscope ? (
          <motion.div
            key="horoscope"
            initial={{ opacity: 0, y: -20, clipPath: "inset(0 0 100% 0)" }}
            animate={{
              opacity: 1,
              y: 0,
              clipPath: "inset(0 0 0% 0)",
            }}
            exit={{ opacity: 0, y: -20 }}
            transition={{
              duration: 1.0,
              ease: [0.4, 0, 0.2, 1],
              clipPath: { duration: 1.0, ease: [0.4, 0, 0.2, 1] },
            }}
            className="flex max-w-3xl px-8 lg:px-2 text-white text-lg/8 text-pretty text-shadow-md"
          >
            <div className="flex flex-col">
              <p>{horoscope}</p>

              {feedbackVisible &&
                <div className="flex flex-col mt-10 items-center gap-4 p-2">
                  <motion.div
                    className="flex mt-10 justify-center"
                    initial={{ opacity: 0, scale: 1 }}
                    animate={{ opacity: 1, scale: 1.05 }}
                    layout
                    transition={{ duration: 3, ease: easeInOut }}
                  >
                    <div className="flex flex-col  p-2">
                      <Feedback paymentIntentId={paymentIntentId} horoscopeId={generateHoroscope?.data?.data?.external_id} />


                    </div>
                  </motion.div>


                  <motion.div
                    className="flex mt-10 justify-center items-center"
                    initial={{ opacity: 0, scale: 1 }}
                    animate={{ opacity: 1, scale: 1.05 }}
                    transition={{ duration: 3, ease: easeInOut }}
                  >
                    {shareHoroscope.isSuccess ? (
                      <div className="flex flex-col items-center gap-4">
                        <div className="text-center">
                          <p className="font-bold">Your horoscope is now public!</p>
                          <p className="text-sm text-slate-300">Share this link with your friends</p>
                        </div>
                        <div className="relative">
                          <button
                            title={copied ? "Copied!" : "Click to copy"}
                            onClick={handleCopy}
                            className="bg-pink-500 hover:bg-pink-600 text-white rounded-md transition-colors duration-200 cursor-pointer flex items-center overflow-hidden"
                          >
                            <span className="px-4 py-3 text-sm truncate max-w-xs">
                              {shareableLink}
                            </span>
                            <div className="w-px self-stretch bg-white/30" />
                            <div className="px-4 py-3 relative">
                              <AnimatePresence mode="wait" initial={false}>
                                {copied ? (
                                  <motion.div
                                    key="check"
                                    initial={{ scale: 0, opacity: 0 }}
                                    animate={{ scale: 1, opacity: 1 }}
                                    exit={{ scale: 0, opacity: 0 }}
                                    transition={{ duration: 0.2 }}
                                  >
                                    <Check size={20} />
                                  </motion.div>
                                ) : (
                                  <motion.div
                                    key="copy"
                                    initial={{ scale: 0, opacity: 0 }}
                                    animate={{ scale: 1, opacity: 1 }}
                                    exit={{ scale: 0, opacity: 0 }}
                                    transition={{ duration: 0.2 }}
                                  >
                                    <Copy size={20} />
                                  </motion.div>
                                )}
                              </AnimatePresence>
                            </div>
                          </button>
                          <AnimatePresence>
                            {copied && (
                              <motion.div
                                key="tooltip"
                                initial={{ opacity: 0, y: 4 }}
                                animate={{ opacity: 1, y: 0 }}
                                exit={{ opacity: 0, y: -4 }}
                                className="absolute -top-9 left-1/2 -translate-x-1/2 bg-slate-800 text-white text-xs px-2 py-1 rounded whitespace-nowrap"
                              >
                                Copied to clipboard!
                              </motion.div>
                            )}
                          </AnimatePresence>
                        </div>
                        <a
                          href={shareableLink}
                          target="_blank"
                          rel="noopener"
                          className="text-sm text-slate-300 hover:text-white underline"
                        >
                          Open in new tab →
                        </a>
                      </div>
                    ) : (

                      <button onClick={() => {
                        // Handle sharing
                        shareHoroscope.mutate({
                          // This exists at this point
                          horoscopeId: generateHoroscope?.data?.data?.external_id!,
                          paymentIntentId: paymentIntentId!,
                        })
                      }} className="bg-pink-500 p-4 text-white rounded-md hover:bg-pink-600 transition duration-200 ease-in-out hover:cursor-pointer">
                        <div className="flex items-center justify-center gap-2">
                          <p>Share Your Horoscope</p>
                          <SquareArrowOutUpRight size={20} />
                        </div>
                      </button>
                    )}
                  </motion.div>


                </div>
              }


            </div>

          </motion.div>
        ) : (
          <motion.form
            key="form"
            initial={{ opacity: 1, y: 0 }}
            exit={{
              opacity: 0,
            }}
            transition={{
              duration: 0.5,
              ease: [0.4, 0, 0.2, 1],
            }}
            id="payment-form"
            onSubmit={handleSubmit}
            className="items-center justify-center flex flex-col w-full px-8 lg:px-4 min-h-[400px]"
          >
            <PaymentElement
              id="payment-element"
              options={paymentElementOptions}
              className="w-full"
            />

            {errorMessage ? (
              <PaymentIntentError message={errorMessage} />
            ) : null}
            <Disclaimer
              onCheckedChange={(checked) => {
                setIsDisclaimerChecked(checked);
                if (checked) {
                  setErrorMessage("");
                }
              }}
            />

            <button
              type="submit"
              disabled={isLoading || !stripe || !elements}
              className={`inline-flex  items-center gap-x-2 rounded-md ${button.background} px-2 py-2.5 text-lg font-semibold text-white shadow-xs transition-colors duration-200 ${button.backgroundHover} ${button.cursor} focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-pink-600`}
            >
              <span id="button-text">
                {isLoading ? (
                  <div className="flex items-center gap-x-2">
                    <span>Consulting the cosmos...</span>
                    <Spinner />
                  </div>
                ) : (
                  "Awesome, let's go!"
                )}
              </span>
            </button>
          </motion.form>
        )}
      </AnimatePresence>
    </div>
  );
}

type GetButtonColorsResponse = {
  cursor: string;
  background: string;
  backgroundHover: string;
};
