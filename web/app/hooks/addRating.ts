import { useMutation } from "@tanstack/react-query";
import axios from "axios";

const API_URL = import.meta.env.VITE_API_URL;
type AddRatingResponse = {
    message: string;
};

export type AddRatingRequest = {
    horoscopeId: string
    paymentIntentId: string
    rating: "positive" | "negative" | "neutral"
}

export const useAddRating = () => {
    return useMutation({
        mutationFn: ({ horoscopeId, paymentIntentId, rating }: AddRatingRequest) =>
            axios.patch<AddRatingResponse>(`${API_URL}/api/v1/horoscopes/${horoscopeId}/rate`, {
                payment_intent_id: paymentIntentId,
                rating
            }),
    });
};
