import { useMutation } from "@tanstack/react-query";
import axios from "axios";

const API_URL = import.meta.env.VITE_API_URL;
type AddRatingResponse = {
    message: string;
};

type AddRatingRequest = {
    horoscopeId: string
    paymentIntentId: string
    rating: string
}

export const useAddRating = ({ horoscopeId, paymentIntentId, rating }: AddRatingRequest) => {
    return useMutation({
        mutationFn: () =>
            axios.patch<AddRatingResponse>(`${API_URL}/api/v1/horoscopes/${horoscopeId}/rate`, {
                payment_intent_id: paymentIntentId,
                rating
            }),
    });
};
