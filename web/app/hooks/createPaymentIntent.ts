import { useMutation } from "@tanstack/react-query";
import axios from "axios";

const API_URL = import.meta.env.VITE_API_URL;
type CreatePaymentIntentResponse = {
  message: string;
  client_secret: string;
  payment_intent_id: string;
};

export const useCreatePaymentIntent = () => {
  return useMutation({
    mutationFn: () =>
      axios.post<CreatePaymentIntentResponse>(`${API_URL}/api/v1/payment-intents`),
  });
};
