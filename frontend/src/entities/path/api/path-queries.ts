import { queryOptions } from "@tanstack/react-query";

import { request } from "@/shared/api";

import type { Path } from "../model/types";

export const pathQueries = {
  all: () =>
    queryOptions({
      queryKey: ["path"],
      queryFn: () => request<Path>("/path"),
    }),
};
