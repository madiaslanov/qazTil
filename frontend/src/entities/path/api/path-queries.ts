import { queryOptions } from "@tanstack/react-query";

import { request } from "@/shared/api";
import { USE_MOCKS } from "@/shared/config/env";
import type { Path } from "../model/types";
import { mockPath } from "./mock";

export const pathQueries = {
  all: () =>
    queryOptions({
      queryKey: ["path"],
      queryFn: () => (USE_MOCKS ? Promise.resolve(mockPath) : request<Path>("/path")),
    }),
};
