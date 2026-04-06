export function PaymentIntentError({ message }: { message: string }) {
  return (
    <div className="my-2 max-w-2xl text-center text-balance">
      <p className="text-sm text-rose-400">{message}</p>
    </div>
  );
}
