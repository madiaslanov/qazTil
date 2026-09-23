import { request } from "@/shared/api";
import type { Progress } from "../model/types";

export const progressApi = {
  list: () => request<Progress[]>("/progress"),
};
