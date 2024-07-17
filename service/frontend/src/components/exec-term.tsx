"use client";

import dynamic from "next/dynamic";

const Terminal = dynamic(() => import("@/components/terminal"), {
  ssr: false,
});

type ExecTermProps = {
  className?: string;
  id: string;
  path: string;
};

const ExecTerm: React.FC<ExecTermProps> = (props) => {
  const { className, id, path } = props;

  return (
    <div id={id} key={id} className={className}>
      <Terminal id="terminal" path={path} className="w-full h-full" />
    </div>
  );
};

export default ExecTerm;
