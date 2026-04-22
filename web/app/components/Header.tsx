type HeaderProps = {
  showTagline?: boolean;
};

export function Header({ showTagline = true }: HeaderProps = {}) {
  return (
    <div className="flex flex-col">
      <div className="flex justify-center">
        <img
          src="/ccblinkie2.gif"
          alt="small banner showing 'credit card horoscope'"
        ></img>
      </div>

      {showTagline && (
        <div className="my-2 text-center">
          <p className="text-lg text-white font-bold">
            What does your credit card say about you?
          </p>
        </div>
      )}
    </div>
  );
}
