"use client";

import dynamic from "next/dynamic";
import FileTree from "@/components/file-tree";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";
import { useEffect, useState } from "react";
import ExecTerm from "@/components/exec-term";
import { useDevenvGeneration } from "@/hooks/use-devenv-generation";
import { useDevenvGenerationMutation } from "@/hooks/use-devenv-generation-mutation";

const Editor = dynamic(() => import("@/components/editor"), {
  ssr: false,
});

export default function Page({ params }: { params: { slug: string } }) {
  const [currentFile, setCurrentFile] = useState<string>();

  const devenvGenerationMutation = useDevenvGenerationMutation({
    uuid: params.slug,
  });
  const generation = useDevenvGeneration({ uuid: params.slug });

  useEffect(() => {
    return () => {
      devenvGenerationMutation.mutate(null);
    };
  }, []);

  return (
    <main className="h-screen w-screen pt-16">
      <ResizablePanelGroup direction="horizontal">
        <ResizablePanel id="file-tree-panel" defaultSize={15}>
          <div className="flex flex-col w-full h-full p-4 space-y-5">
            <FileTree
              className="flex flex-col w-full h-full space-y-5"
              devenvUuid={params.slug}
              selectedFile={currentFile}
              setSelectedFile={setCurrentFile}
            />
          </div>
        </ResizablePanel>
        <ResizableHandle />
        <ResizablePanel defaultSize={85}>
          <ResizablePanelGroup direction="vertical">
            <ResizablePanel id="editor-panel">
              <Editor
                className="w-full h-full"
                devenvUuid={params.slug}
                filename={currentFile}
              />
            </ResizablePanel>
            {Boolean(generation) && (
              <>
                <ResizableHandle />
                <ResizablePanel id="terminal-panel">
                  <ExecTerm
                    className="w-full h-full"
                    id={String(generation)}
                    path={"/api/devenv/" + params.slug + "/exec"}
                  />
                </ResizablePanel>
              </>
            )}
          </ResizablePanelGroup>
        </ResizablePanel>
      </ResizablePanelGroup>
    </main>
  );
}
