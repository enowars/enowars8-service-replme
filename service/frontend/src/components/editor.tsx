"use client";

import { useDevenvFileContentQuery } from '@/hooks/use-devenv-file-content-query';
import { useDevenvFileContentMutation } from '@/hooks/use-devenv-file-content-mutation';
import MonacoEditor, { Monaco, OnChange } from '@monaco-editor/react';
import { useTheme } from 'next-themes';

type EditorProps = {
  className?: string,
  devenvUuid: string,
  filename?: string
}

const Editor: React.FC<EditorProps> = (props) => {
  const { className, devenvUuid, filename } = props;

  const { resolvedTheme } = useTheme();
  const editorTheme = resolvedTheme === "light" ? "light" : "dark";

  const fileContentQuery = useDevenvFileContentQuery({
    uuid: devenvUuid,
    filename
  })

  const fileContentMutation = useDevenvFileContentMutation({
    uuid: devenvUuid,
    filename
  })

  const handleEditorWillMount = (monaco: Monaco) => {
    monaco.editor.defineTheme("dark", {
      "base": "vs-dark",
      "inherit": true,
      "rules": [],
      "colors": {
        "editor.background": "#020817"
      }
    })
  }

  const handleEditorChange: OnChange = (value, _) => {
    if (value)
      fileContentMutation.mutate(value);
  }

  if (!filename)
    return <></>

  if (fileContentQuery.isStale || fileContentQuery.isLoading) {
    return <></>
  }

  return (
    <MonacoEditor
      key={filename}
      className={className}
      defaultValue={fileContentQuery.data}
      defaultLanguage='c'
      theme={editorTheme}
      beforeMount={handleEditorWillMount}
      onChange={handleEditorChange}
    />
  )
}

export default Editor;
