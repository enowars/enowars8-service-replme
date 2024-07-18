"use client";
import Image from "next/image";
import { useEffect, useState } from "react";

const IMAGES = [
  "/assets/imagi-date.png",
  "/assets/onlyflags.png",
  "/assets/onlyflags2.png",
  "/assets/pirate-say.png",
  "/assets/sceam.png",
  "/assets/whatsscam.png",
];

function getRandomArbitrary(min: number, max: number) {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

const RandomAdBanner = () => {
  const [index, setIndex] = useState<number>();

  useEffect(() => {
    setIndex(getRandomArbitrary(0, IMAGES.length - 1));
  }, []);

  if (index === undefined) return <></>;

  return (
    <Image
      className="object-contain"
      src={IMAGES[index]}
      fill={true}
      alt="ad"
    />
  );
};

export default RandomAdBanner;
