import { Link } from "react-router";

type DisclaimerProps = {
  onCheckedChange?: (checked: boolean) => void;
};

export function Disclaimer({ onCheckedChange }: DisclaimerProps) {
  return (
    <label htmlFor="offers" className="flex gap-3 py-4 cursor-pointer">
      <div className="flex h-6 shrink-0 items-center">
        <div className="group grid size-4 grid-cols-1">
          <input
            id="offers"
            name="offers"
            type="checkbox"
            required
            aria-describedby="offers-description"
            onChange={(e) => {
              onCheckedChange?.(e.target.checked);
            }}
            className="col-start-1 row-start-1 appearance-none rounded-sm border border-gray-300 bg-white checked:border-pink-600 checked:bg-pink-600 indeterminate:border-pink-600 indeterminate:bg-pink-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-pink-600 disabled:border-gray-300 disabled:bg-gray-100 disabled:checked:bg-gray-100 forced-colors:appearance-auto cursor-pointer"
          />
          <svg
            fill="none"
            viewBox="0 0 14 14"
            className="pointer-events-none col-start-1 row-start-1 size-3.5 self-center justify-self-center stroke-white group-has-disabled:stroke-gray-950/25"
          >
            <path
              d="M3 8L6 11L11 3.5"
              strokeWidth={2}
              strokeLinecap="round"
              strokeLinejoin="round"
              className="opacity-0 group-has-checked:opacity-100"
            />
            <path
              d="M3 7H11"
              strokeWidth={2}
              strokeLinecap="round"
              strokeLinejoin="round"
              className="opacity-0 group-has-indeterminate:opacity-100"
            />
          </svg>
        </div>
      </div>
      <div className="text-sm/6">
        <span id="offers-description" className="text-slate-200">
          I understand this AI horoscope is for entertainment only and agree to
          the{" "}
          <Link
            to="/terms-of-service"
            onClick={(e) => e.stopPropagation()}
            className="text-blue-500 hover:text-blue-700 transition duration-200 hover:underline"
          >
            Terms of Service
          </Link>{" "}
          and{" "}
          <Link
            to="/privacy-policy"
            onClick={(e) => e.stopPropagation()}
            className="text-blue-500 hover:text-blue-700 transition duration-200 hover:underline"
          >
            Privacy Policy
          </Link>
          .
        </span>
      </div>
    </label>
  );
}
