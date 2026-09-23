import { request } from "@/shared/api";
import type { Category } from "../model/types";

export const categoryApi = {
  list: () => request<Category[]>("/categories"),
  get: (id: number) => request<Category>(`/categories/${id}`),
};
