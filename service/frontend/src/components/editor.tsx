"use client";

import { useDevenvFileContentQuery } from "@/hooks/use-devenv-file-content-query";
import { useDevenvFileContentMutation } from "@/hooks/use-devenv-file-content-mutation";
import MonacoEditor, { Monaco, OnChange } from "@monaco-editor/react";
import { useTheme } from "next-themes";
import { useEffect, useState } from "react";
import { useDebounce } from "@uidotdev/usehooks";
import { useDevenvStateMutation } from "@/hooks/use-devenv-state-mutation";

type EditorProps = {
  className?: string;
  devenvUuid: string;
  filename?: string;
};

type DebouncedEditorProps = {
  className?: string;
  devenvUuid: string;
  filename: string;
  initialContent: string;
};

const DebouncedEditor: React.FC<DebouncedEditorProps> = (props) => {
  const { className, devenvUuid, filename, initialContent } = props;

  const { resolvedTheme } = useTheme();
  const editorTheme = resolvedTheme === "light" ? "light" : "dark";

  const [content, setContent] = useState<string>(initialContent);
  const debouncedContent = useDebounce(content, 1000);

  const fileContentMutation = useDevenvFileContentMutation({
    uuid: devenvUuid,
    filename,
  });

  const devenvStateMutation = useDevenvStateMutation({
    uuid: devenvUuid,
  });

  useEffect(() => {
    fileContentMutation.mutate(debouncedContent);
  }, [debouncedContent, fileContentMutation.mutate]);

  const handleEditorChange: OnChange = (value, _) => {
    if (value) {
      devenvStateMutation.mutate("dirty");
      setContent(value);
    }
  };

  useEffect(() => {
    return () => {
      devenvStateMutation.mutate("ok");
    };
  }, []);

  const handleEditorWillMount = (monaco: Monaco) => {
    monaco.editor.defineTheme("dark", {
      base: "vs-dark",
      inherit: true,
      rules: [],
      colors: {
        "editor.background": "#020817",
      },
    });
  };

  return (
    <MonacoEditor
      key={filename}
      className={className}
      defaultValue={initialContent}
      defaultLanguage="c"
      theme={editorTheme}
      beforeMount={handleEditorWillMount}
      onChange={handleEditorChange}
    />
  );
};

const Editor: React.FC<EditorProps> = (props) => {
  const { className, devenvUuid, filename } = props;

  const fileContentQuery = useDevenvFileContentQuery({
    uuid: devenvUuid,
    filename,
  });

  if (!filename) return <></>;

  if (fileContentQuery.isStale || fileContentQuery.isLoading) {
    return <></>;
  }

  if (!fileContentQuery.isSuccess) {
    return <>Loading file failed</>;
  }

  return (
    <DebouncedEditor
      className={className}
      devenvUuid={devenvUuid}
      filename={filename}
      initialContent={fileContentQuery.data}
    />
  );
};

export default Editor;
