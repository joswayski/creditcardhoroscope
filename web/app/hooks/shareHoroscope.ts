import { useMutation } from "@tanstack/react-query";
import axios from "axios";

const API_URL = import.meta.env.VITE_API_URL;
type ShareHoroscopeBody = {
    paymentIntentId: string;
    horoscopeId: string;
};

type ShareHoroscopeResponse = {
    message: string;
    //   horoscope: string;
    //   external_id?: string;
};

export const useShareHoroscope = () => {
    return useMutation({
        mutationFn: ({ paymentIntentId, horoscopeId }: ShareHoroscopeBody) =>
            axios.post<ShareHoroscopeResponse>(`${API_URL}/api/v1/horoscopes/${horoscopeId}/share`, {
                payment_intent_id: paymentIntentId,
            }),
    });
};
