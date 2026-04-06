import { Link } from "react-router";

const currentYear = new Date().getFullYear();

export const Footer = () => {
  return (
    <footer className="text-sm text-slate-500 flex flex-col sm:flex-row justify-center items-center py-4 gap-2 px-4">
      <span>© {currentYear} Credit Card Horoscope</span>
      <span className="hidden sm:inline">|</span>
      <div className="flex items-center gap-x-2">
        <Link
          to="/terms-of-service"
          className="hover:text-slate-700 hover:underline transition duration-200"
        >
          Terms of Service
        </Link>
        <span>|</span>
        <Link
          to="/privacy-policy"
          className="hover:text-slate-700 hover:underline transition duration-200"
        >
          Privacy Policy
        </Link>
      </div>
    </footer>
  );
};
