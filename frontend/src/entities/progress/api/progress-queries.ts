import { queryOptions } from "@tanstack/react-query";

import { request } from "@/shared/api";
import type { Progress } from "../model/types";

export const progressQueries = {
  all: () =>
    queryOptions({
      queryKey: ["progress"],
      queryFn: () => request<Progress[]>("/progress"),
    }),
};
