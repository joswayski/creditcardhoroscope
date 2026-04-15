import type { Route } from "./+types/home";
import { FindOutButton } from "~/components/FindOutButton";
import { loadStripe, type Appearance } from "@stripe/stripe-js";
import { Header } from "~/components/Header";
import { useCreatePaymentIntent } from "~/hooks/createPaymentIntent";
import { Elements } from "@stripe/react-stripe-js";
import { CheckoutForm } from "~/components/CheckoutForm";
import { Footer } from "~/components/Footer";

const stripePromise = loadStripe(import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY);

export function meta({ }: Route.MetaArgs) {
  return [
    { title: "Credit Card Horoscope" },
    {
      name: "description",
      content: "What does your credit card say about you?",
    },
  ];
}

const appearance: Appearance = {
  theme: "stripe",
};
const loader = "auto";

export default function Home() {
  const createPi = useCreatePaymentIntent();
  const clientSecret = createPi?.data?.data?.client_secret;
  return (
    <div className=" flex flex-col min-h-screen">
      <main className="flex-1 items-center justify-center pt-16 pb-4 flex-col">
        <Header />

        <div className="my-4  justify-center flex ">
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
            <FindOutButton createPi={createPi} />
          )}
        </div>
      </main>
      <div
        className="sticky bottom-0 z-10"
        style={{
          backgroundImage: "url('/stars3.gif')",
          backgroundRepeat: "repeat",
        }}
      >
        <Footer />
      </div>
    </div>
  );
}
