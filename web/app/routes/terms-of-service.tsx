import { BiArrowBack } from "react-icons/bi";
import type { Route } from "../+types/root";
import { Link } from "react-router";

export function meta({ }: Route.MetaArgs) {
  return [
    { title: "Credit Card Horoscope Terms of Service" },
    {
      name: "description",
      content: "View the Terms of Service for Credit Card Horoscope",
    },
  ];
}

export default function TOS() {
  return (
    <div className="flex flex-col ">
      <div className="lg:px-20 px-8 pt-8">
        <div className="flex flex-row justify-start items-center  max-w-6xl">
          <Link
            to="/"
            className="inline-flex items-center gap-x-1.5 rounded-md bg-pink-600 px-3 py-2 text-sm font-semibold text-white shadow-xs hover:bg-pink-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-pink-600 transition duration-200"
          >
            <BiArrowBack aria-hidden="true" className="-ml-0.5 size-5" />
            Home
          </Link>
        </div>
      </div>
      <div className="flex flex-col lg:p-20 p-8 justify-center items-center ">
        <div className="p-4  max-w-6xl space-y-4 text-pretty">
          <h1 className="text-4xl font-extrabold">
            Terms of Service for creditcardhoroscope.com
          </h1>
          <p className="font-bold text-2xl">Last updated: November 15, 2025</p>

          <div className="py-4 space-y-8">
            <p>
              Welcome to{" "}
              <span className="font-bold">creditcardhoroscope.com</span>. By
              using our website and paying for our service (the "Service"), you
              are agreeing to these Terms of Service (“Terms”). Please read them
              carefully. This is a legally binding agreement between you and{" "}
              <span className="font-bold">Valerio Group LLC</span> ("we", "us",
              "our").
            </p>

            <div>
              <h4 className="font-bold text-xl">
                1. Description of Our Service
              </h4>
              <p>
                <span className="font-bold">creditcardhoroscope.com</span>{" "}
                provides an AI-generated "horoscope" (the "Horoscope") for a
                one-time fee of $1.00 USD. The Horoscope is generated based on
                non-sensitive payment metadata.
              </p>
            </div>


            <div>
              <h4 className="font-bold text-xl">
                2. Entertainment Purposes Only
              </h4>
              <p>You expressly understand and agree to the following:</p>
              <ul className="list-disc list-inside ml-4 space-y-2">
                <li>
                  The Service provides a fictional, AI-generated horoscope for
                  entertainment purposes only.
                </li>
                <li>
                  The Horoscope is not real, not factual, and not based on any
                  science.
                </li>
                <li>
                  The Horoscope is not financial advice, legal advice, or
                  personal advice of any kind.
                </li>
                <li>
                  You agree not to use the Horoscope to make any life,
                  financial, or personal decisions.
                </li>
              </ul>
            </div>

            <div>
              <h4 className="font-bold text-xl">
                3. No Guarantees or Warranties
              </h4>
              <p>
                The Service is provided "as is" and "as available" without any
                warranties. We do not guarantee the accuracy, completeness, or
                reliability of any AI-generated Horoscope. We do not guarantee
                that the Service will be available, uninterrupted, or
                error-free.
              </p>
            </div>

            <div>
              <h4 className="font-bold text-xl">
                4. Payment and Refund Policy
              </h4>
              <ul className="list-disc list-inside ml-4 space-y-2">
                <li>
                  <span className="font-bold">Payment</span>: All payments are
                  processed securely by our third-party payment processor,
                  Stripe, Inc. Here are links to their{" "}
                  <a
                    href="https://stripe.com/legal/ssa"
                    target="_blank"
                    className="inline-block text-blue-500 hover:text-blue-700 transition duration-200 hover:underline px-1 -mx-1 cursor-pointer"
                  >
                    Terms of Service
                  </a>{" "}
                  and{" "}
                  <a
                    href="https://stripe.com/privacy"
                    target="_blank"
                    className="inline-block text-blue-500 hover:text-blue-700 transition duration-200 hover:underline px-1 -mx-1 cursor-pointer"
                  >
                    Privacy Policy
                  </a>{" "}
                  .
                </li>
                <li>

                  <span className="font-bold">Refunds</span>: We want you to be
                  satisfied. If you are unhappy with your Horoscope, please
                  contact us at{" "}
                  <span className="font-bold">contact@josevalerio.com</span> for
                  a full $1.00 refund. We would rather give you a refund than
                  have you file a dispute.
                </li>
              </ul>
            </div>
            <div>
              <h4 className="font-bold text-xl">5. Limitation of Liability</h4>
              <p>
                You agree that Valerio Group LLC is not liable to you or any
                third party for any damages of any kind that result from the use
                of the Service. Our maximum liability to you for any reason
                shall be limited to the amount you paid for the Service (i.e.,
                $1.00 USD).
              </p>
            </div>

            <div>
              <h4 className="font-bold text-xl">6. Intellectual Property</h4>
              <p>
                We own all rights, title, and interest in and to the
                creditcardhoroscope.com website, including all design, text, and
                code. You may not duplicate, copy, or reuse any portion of our
                website without our express written permission.
              </p>
            </div>

            <div>
              <h4 className="font-bold">7. Changes to These Terms</h4>
              <p>
                We may update these Terms at any time. If we do, we will refresh
                the "Last updated" date at the top of this page. Your continued
                use of the Service after any changes constitutes your acceptance
                of the new Terms.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
