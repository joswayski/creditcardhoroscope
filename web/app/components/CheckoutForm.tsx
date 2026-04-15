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
import { Frown, Meh, Smile } from 'lucide-react';
import { useAddRating, type AddRatingRequest } from "~/hooks/addRating";


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

  if (!addRating.isIdle) {
    return <p className="text-center font-bold">Thanks for your feedback!</p>
  }

  return <div className="flex max-w-sm justify-center">
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
  </div >
}


export function CheckoutForm({ clientSecret }: CheckoutFormProps) {
  const stripe = useStripe();
  const elements = useElements();
  const generateHoroscope = useGenerateHoroscope();
  const [errorMessage, setErrorMessage] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [isDisclaimerChecked, setIsDisclaimerChecked] = useState(false);
  const [feedbackVisible, setFeedbackVisible] = useState(false)
  const [paymentIntentId, setPaymentIntentId] = useState<null | string>(null)
  const [horoscopeId, setHoroscopeId] = useState<null | string>(null)

  const isButtonDisabled =
    isLoading || !stripe || !elements || !isDisclaimerChecked;
  const button = getButtonColors({
    isLoading,
    isDisabled: !isDisclaimerChecked,
  });

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
      }, 8000)

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
                <motion.div
                  className="flex mt-10 justify-center"
                  initial={{ opacity: 0, scale: 1 }}
                  animate={{ opacity: 1, scale: 1.05 }}
                  layout
                  transition={{ duration: 3, ease: easeInOut }}
                >
                  <Feedback paymentIntentId={paymentIntentId} horoscopeId={generateHoroscope?.data?.data?.external_id} />
                </motion.div>
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
