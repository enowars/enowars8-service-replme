"use client";

import dynamic from 'next/dynamic'

const Terminal = dynamic(() => import('@/components/terminal'), {
  ssr: false
})

export default function Page({ params }: { params: { slug: string } }) {
  return (
    <div className="w-screen h-screen pt-16 px-2 pb-2" >
      <Terminal id="terminal" name={params.slug} className="w-full h-full" />
    </div>
  )
}
