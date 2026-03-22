import type { Route } from "./+types/home";
import { Welcome } from "../welcome/welcome";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Credit Card Horoscope" },
    {
      name: "description",
      content: "What does your credit card say about you?",
    },
  ];
}

export default function Home() {
  return <Welcome />;
}
