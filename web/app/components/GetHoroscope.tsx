import { loadStripe, type Appearance } from "@stripe/stripe-js";
import { Elements } from "@stripe/react-stripe-js";
import { FindOutButton } from "./FindOutButton";
import { CheckoutForm } from "./CheckoutForm";
import { useCreatePaymentIntent } from "~/hooks/createPaymentIntent";

const stripePromise = loadStripe(import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY);

const appearance: Appearance = {
  theme: "stripe",
};
const loader = "auto";

type GetHoroscopeProps = {
  cta?: string;
};

export function GetHoroscope({ cta }: GetHoroscopeProps = {}) {
  const createPi = useCreatePaymentIntent();
  const clientSecret = createPi?.data?.data?.client_secret;

  return (
    <div className="my-4 justify-center flex">
      {createPi.isSuccess && clientSecret ? (
        <div>
          <Elements
            options={{
              clientSecret,
              appearance,
              loader,
            }}
            stripe={stripePromise}
          >
            <CheckoutForm clientSecret={clientSecret} />
          </Elements>
        </div>
      ) : (
        <FindOutButton createPi={createPi} cta={cta} />
      )}
    </div>
  );
}
