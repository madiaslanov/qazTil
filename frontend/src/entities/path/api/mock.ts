import type { Path } from "../model/types";

export const mockPath: Path = {
  units: [
    {
      id: "1",
      title_kk: "Дәмхана",
      title_ru: "Кафе",
      description: "Заказ и счёт",
      level: "A0",
      topic: "life",
      position: 1,
      lessons: [
        {
          id: "1",
          title_kk: "Сәлемдесу",
          title_ru: "Приветствие",
          position: 1,
          state: "completed",
          correct_answers: 5,
          total_answers: 5,
        },
        {
          id: "2",
          title_kk: "Тапсырыс",
          title_ru: "Заказ",
          position: 2,
          state: "available",
          resume_position: 2,
          correct_answers: 0,
          total_answers: 0,
        },
        {
          id: "3",
          title_kk: "Есеп",
          title_ru: "Счёт",
          position: 3,
          state: "locked",
          correct_answers: 0,
          total_answers: 0,
        },
      ],
    },
  ],
};
