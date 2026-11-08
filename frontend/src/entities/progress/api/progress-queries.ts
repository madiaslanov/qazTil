import { queryOptions } from "@tanstack/react-query";

import { request } from "@/shared/api";
import type { CourseProgress } from "../model/types";

export const progressQueries = {
  all: () =>
    queryOptions({
      queryKey: ["progress"],
      queryFn: () => request<CourseProgress>("/progress"),
    }),
};
