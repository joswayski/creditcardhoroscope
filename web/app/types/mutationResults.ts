import type { UseMutationResult } from "@tanstack/react-query";
import type { AxiosResponse } from "axios";

export type MutationResult = UseMutationResult<AxiosResponse, Error, void>;
