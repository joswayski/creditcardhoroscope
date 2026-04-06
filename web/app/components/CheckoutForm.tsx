import {
  PaymentElement,
  useElements,
  useStripe,
} from "@stripe/react-stripe-js";
import type { StripePaymentElementOptions } from "@stripe/stripe-js";
import { useState } from "react";
import { Disclaimer } from "./Disclaimer";
import { PaymentIntentError } from "./PiError";
import { useGenerateHoroscope } from "~/hooks/generateHoroscope";
import { Spinner } from "./Spinner";
import { AnimatePresence, motion } from "motion/react";
import axios, { AxiosError } from "axios";

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

export function CheckoutForm({ clientSecret }: CheckoutFormProps) {
  const stripe = useStripe();
  const elements = useElements();
  const generateHoroscope = useGenerateHoroscope();
  const [errorMessage, setErrorMessage] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [isDisclaimerChecked, setIsDisclaimerChecked] = useState(false);
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

  return (
    <div className="relative min-h-[400px] overflow-hidden">
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
            <p>{horoscope}</p>
          </motion.div>
        ) : (
          <motion.form
            key="form"
            initial={{ opacity: 1, y: 0 }}
            exit={{
              opacity: 0,
              y: 20,
              clipPath: "inset(100% 0 0 0)",
            }}
            transition={{
              duration: 1.0,
              ease: [0.4, 0, 0.2, 1],
              clipPath: { duration: 1.0, ease: [0.4, 0, 0.2, 1] },
            }}
            id="payment-form"
            onSubmit={handleSubmit}
            className="items-center justify-center flex flex-col w-full px-8 lg:px-4"
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
