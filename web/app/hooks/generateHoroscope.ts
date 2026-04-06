import { useMutation } from "@tanstack/react-query";
import axios from "axios";

const API_URL = import.meta.env.VITE_API_URL;
type GenerateHoroscopeBody = {
  paymentIntentId: string;
};

type GenerateHoroscopeResponse = {
  message: string;
  horoscope: string;
};

export const useGenerateHoroscope = () => {
  return useMutation({
    mutationFn: ({ paymentIntentId }: GenerateHoroscopeBody) =>
      axios.post<GenerateHoroscopeResponse>(`${API_URL}/api/v1/horoscopes`, {
        paymentIntentId,
      }),
  });
};
