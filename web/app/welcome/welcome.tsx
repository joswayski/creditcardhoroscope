export function Welcome() {
  return (
    <main className="flex items-center justify-center pt-16 pb-4 ">
      <div className="flex-1 flex flex-col items-center gap-4 min-h-0 text-pink-700">
        <h1 className="text-5xl font-bold">Credit Card Horoscope</h1>
        <p>Coming soon! Follow along:</p>
        <a
          target="_blank"
          rel="noopener noreferrer"
          className="text-blue-400 hover:text-blue-600 transition duration 200ms ease-in-out underline"
          href="https://github.com/joswayski/creditcardhoroscope"
        >
          https://github.com/joswayski/creditcardhoroscope
        </a>
      </div>
    </main>
  );
}
