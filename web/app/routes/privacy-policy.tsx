import { BiArrowBack } from "react-icons/bi";
import type { Route } from "../+types/root";
import { Link } from "react-router";

export function meta({ }: Route.MetaArgs) {
  return [
    { title: "Credit Card Horoscope Privacy Policy" },
    {
      name: "description",
      content: "View the Privacy Policy for Credit Card Horoscope",
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
          <h1 className="text-4xl font-extrabold">Privacy Policy for creditcardhoroscope.com</h1>
          <p className="font-bold text-2xl">Last updated: November 15, 2025</p>

          <div className="py-4 space-y-8">
            <p>
              This policy explains what data we collect, why we collect it, and what we do with it.
            </p>
            <div>
              <h4 className="font-bold text-xl">1. What We Collect and Why</h4>
              <p>
                We use third-party services to collect the necessary data to provide our Service.
              </p>
              <ul className="list-disc list-inside ml-4 space-y-2">
                <li>
                  <span className="font-bold">Security & Analytics Data</span>: We use{" "}
                  <a
                    href="https://cloudflare.com/"
                    target="_blank"
                    className="inline-block text-blue-500 hover:text-blue-700 transition duration-200 hover:underline px-1 -mx-1 cursor-pointer"
                  >
                    Cloudflare
                  </a>{" "}
                  as a security and performance layer. To do this, Cloudflare processes user data on
                  our behalf, which may include your IP Address, browser type, and other details
                  used to detect malicious traffic and provide basic analytics.
                </li>

                <li>
                  <span className="font-bold">Payment Information</span>: We use{" "}
                  <a
                    href="https://stripe.com/"
                    target="_blank"
                    className="inline-block text-blue-500 hover:text-blue-700 transition duration-200 hover:underline px-1 -mx-1 cursor-pointer"
                  >
                    Stripe
                  </a>{" "}
                  as our exclusive, third-party payment processor.
                </li>
                <li>
                  <span className="font-bold">What We DO NOT Handle</span>: We never see, collect,
                  or store your full credit card number or CVC. This information goes directly to
                  Stripe and never touches our servers.
                </li>
                <li>
                  <span className="font-bold"> What We DO Receive and Store</span>: To generate your
                  Horoscope and for our business records, we receive and store the following
                  information:
                </li>
                <ul className="list-disc list-inside ml-4 space-y-2">
                  <li>
                    Your <span className="font-bold">Payment Intent ID</span> from Stripe.
                  </li>
                  <li>
                    The <span className="font-bold">Horoscope(s)</span> we generate for you.
                  </li>
                  <li>
                    The <span className="font-bold">payment metadata</span> provided by Stripe,
                    which includes your card brand, expiration month and year, last 4 digits,
                    country and postal code.
                  </li>
                </ul>
              </ul>
            </div>
            <div>
              <h4 className="font-bold text-xl">2. How We Use Your Information</h4>
              <p>We use the information we collect to:</p>
              <ul className="list-disc list-inside ml-4 space-y-2">
                <li>
                  Secure our website and protect it from malicious attacks (e.g., bots and DDoS
                  attacks)
                </li>

                <li>
                  Generate your AI Horoscope by sending the payment metadata (as described in
                  Section 1) as a unique input to our third-party AI provider.
                </li>
                <li>Process your payment and maintain internal business records.</li>
                <li>
                  Understand our audience and improve our service by analyzing anonymous, aggregated
                  data (e.g., "what percentage of our users are from the US").
                </li>

              </ul>
            </div>
            <div>
              <h4 className="font-bold text-xl">3. Third-Party Services</h4>
              <p>
                We do not sell, rent, or share your personal information with any third parties for
                their marketing purposes. We use the following third-party services to provide our
                Service:
              </p>
              <ul className="list-disc list-inside ml-4 space-y-2">
                <li>
                  <span className="font-bold">Cloudflare, Inc.</span>: To provide security, DDoS
                  protection, and basic analytics.
                </li>

                <li>
                  <span className="font-bold">Stripe, Inc</span>.: To process all payments. Your
                  browser sends your payment information directly to Stripe.
                </li>
                <li>
                  <span className="font-bold">OpenRouter</span>: To generate the AI Horoscope. We
                  send the payment metadata (card brand, expiration, last 4 digits, country, postal
                  code) to this service, and they return the generated Horoscope.
                </li>
              </ul>
            </div>
            <div>
              <h4 className="font-bold text-xl">4. Data Retention</h4>
              <p>
                We retain the information we collect (Payment Intent ID, generated Horoscopes,
                payment metadata
                {/*,  and
                analytics/geolocation data*/}
                ) for our internal business records, to improve our service, and to handle any
                refund requests.
              </p>
            </div>
            <div>
              <h4 className="font-bold text-xl">5. Cookies</h4>
              <p>
                We may use essential cookies to make our site functional (e.g., to process your
                payment) or for basic, anonymous analytics (such as those provided by Cloudflare).
                We do not use third-party tracking or advertising cookies.
              </p>
            </div>

            <div>
              <h4 className="font-bold text-xl">6. Changes and Questions</h4>
              <p>
                We may update this policy. If we do, we will refresh the "Last updated" date at the
                top of this page. If you have any questions about this policy, please contact us at
                contact@josevalerio.com.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
