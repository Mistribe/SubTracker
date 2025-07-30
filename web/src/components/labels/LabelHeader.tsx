interface LabelHeaderProps {
  totalCount: number;
}

export const LabelHeader = ({ totalCount }: LabelHeaderProps) => {
  return (
    <div className="flex justify-between items-center mb-6">
      <h1 className="text-3xl font-bold">Labels</h1>
      <p>{totalCount} elements</p>
    </div>
  );
};