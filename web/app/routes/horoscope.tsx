import type { Route } from "./+types/horoscope";
import axios from "axios";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { Header } from "~/components/Header";
import { Footer } from "~/components/Footer";
import { GetHoroscope } from "~/components/GetHoroscope";

dayjs.extend(relativeTime);

type HoroscopeResponse = {
    id: string;
    horoscope: string;
    created_at: string;
};
export async function loader({ params }: Route.LoaderArgs) {
    const { id } = params;
    try {
        const res = await axios.get<HoroscopeResponse>(
            `${process.env.VITE_API_URL}/api/v1/horoscopes/${id}`
        );
        return { horoscope: res.data };
    } catch (e) {
        throw new Response("Not found", { status: 404 });
    }
}

export function meta({ data, params }: Route.MetaArgs) {
    if (!data?.horoscope) {
        return [{ title: "Horoscope not found" }];
    }

    const preview = data.horoscope.horoscope.slice(0, 160) + "...";
    const url = `https://creditcardhoroscope.com/${params.id}`;

    return [
        { title: "Your Credit Card Horoscope" },
        { name: "description", content: preview },
        { property: "og:title", content: "Your Credit Card Horoscope" },
        { property: "og:description", content: preview },
        { property: "og:url", content: url },
        { property: "og:type", content: "website" },
        { name: "twitter:card", content: "summary" },
        { name: "twitter:title", content: "Your Credit Card Horoscope" },
        { name: "twitter:description", content: preview },
    ];
}

export default function HoroscopePage({ loaderData }: Route.ComponentProps) {
    const { horoscope } = loaderData;

    return (
        <div className="flex flex-col min-h-screen">
            <main className="flex-1 items-center justify-center pt-16 pb-4 flex-col">
                <Header showTagline={false} />
                <div className="mt-8 flex flex-col max-w-3xl px-8 mx-auto text-white text-lg/8 text-pretty text-shadow-md">
                    <p>{horoscope.horoscope}</p>
                    <p className="mt-4 text-sm text-slate-400 text-center">
                        Created {dayjs(horoscope.created_at).fromNow()}
                    </p>
                </div>
                <div className="mt-10 flex justify-center">
                    <GetHoroscope cta="Get yours for $1" />
                </div>
            </main>
            <Footer />
        </div>
    );
}