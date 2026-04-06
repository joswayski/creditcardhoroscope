import { IoSparkles } from "react-icons/io5";
import type { MutationResult } from "~/types/mutationResults";
import { PaymentIntentError } from "./PiError";

type FindOutButtonProps = {
  createPi: MutationResult;
};

type GetButtonColorsResponse = {
  cursor: string;
  background: string;
  backgroundHover: string;
};

const getButtonColors = (mutation: MutationResult): GetButtonColorsResponse => {
  if (mutation.isPending) {
    return {
      cursor: "cursor-not-allowed",
      background: "bg-pink-300",
      backgroundHover: "bg-pink-300",
    };
  }

  return {
    cursor: "cursor-pointer",
    background: "bg-pink-500",
    backgroundHover: "hover:bg-pink-400",
  };
};

export function FindOutButton({ createPi }: FindOutButtonProps) {
  const button = getButtonColors(createPi);

  return createPi.isError ? (
    <PaymentIntentError
      message={`An error ocurred while loading the payment section, please refresh the
        page and try again :(`}
    />
  ) : (
    <button
      onClick={() => createPi.mutate()}
      type="button"
      disabled={createPi.isPending}
      className={`inline-flex  items-center gap-x-2 rounded-md ${button.background} px-3.5 py-2.5 text-lg font-semibold text-white shadow-xs transition-colors duration-200 ${button.backgroundHover} ${button.cursor} focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-pink-600`}
    >
      Find out for $1
      <IoSparkles />
    </button>
  );
}
