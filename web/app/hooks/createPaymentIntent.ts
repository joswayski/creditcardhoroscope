import { useMutation } from "@tanstack/react-query";
import axios from "axios";

const API_URL = import.meta.env.VITE_API_URL;
type CreatePaymentIntentResponse = {
  message: string;
  clientSecret: string;
};

export const useCreatePaymentIntent = () => {
  return useMutation({
    mutationFn: () =>
      axios.post<CreatePaymentIntentResponse>(`${API_URL}/payment-intents`),
  });
};
