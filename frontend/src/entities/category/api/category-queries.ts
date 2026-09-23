import { queryOptions } from "@tanstack/react-query";

import { request } from "@/shared/api";
import type { Category } from "../model/types";

export const categoryQueries = {
  all: () =>
    queryOptions({
      queryKey: ["categories"],
      queryFn: () => request<Category[]>("/categories"),
    }),
};
