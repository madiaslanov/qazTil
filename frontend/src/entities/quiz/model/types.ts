export type Question = {
  id: number;
  prompt: string;
  options: string[];
  selected_index?: number;
  correct?: boolean;
  correct_index?: number;
};

export type Score = {
  correct: number;
  total: number;
};

export type Quiz = {
  id: number;
  category_id: number;
  created_at: string;
  finished_at?: string;
  score: Score;
  questions: Question[];
};

export type AnswerResult = {
  correct: boolean;
  correct_index: number;
  quiz: Quiz;
};
