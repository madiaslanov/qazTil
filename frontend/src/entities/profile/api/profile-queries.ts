import { queryOptions } from "@tanstack/react-query";

import { request } from "@/shared/api";
import { USE_MOCKS } from "@/shared/config/env";
import type { Profile, Stats } from "../model/types";
import { mockProfile, mockStats } from "./mock";

export const profileQueries = {
  profile: () =>
    queryOptions({
      queryKey: ["me", "profile"],
      queryFn: () =>
        USE_MOCKS ? Promise.resolve(mockProfile) : request<Profile>("/me/profile"),
    }),
  stats: () =>
    queryOptions({
      queryKey: ["me", "stats"],
      queryFn: () =>
        USE_MOCKS ? Promise.resolve(mockStats) : request<Stats>("/me/stats"),
    }),
};
