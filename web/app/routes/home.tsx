import type { Route } from "./+types/home";
import { Header } from "~/components/Header";
import { Footer } from "~/components/Footer";
import { GetHoroscope } from "~/components/GetHoroscope";

export function meta({ }: Route.MetaArgs) {
  return [
    { title: "Credit Card Horoscope" },
    {
      name: "description",
      content: "What does your credit card say about you?",
    },
  ];
}

export default function Home() {
  return (
    <div className=" flex flex-col min-h-screen">
      <main className="flex-1 items-center justify-center pt-16 pb-4 flex-col">
        <Header />
        <GetHoroscope />
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
