"use client";

import MonacoEditor, { Monaco, OnChange } from '@monaco-editor/react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import axios from 'axios';
import { useTheme } from 'next-themes';

type EditorProps = {
  className?: string,
  devenvUuid: string,
  filename?: string
}

const Editor: React.FC<EditorProps> = (props) => {
  const { className, devenvUuid, filename } = props;
  const queryClient = useQueryClient();

  const { resolvedTheme } = useTheme();
  const editorTheme = resolvedTheme === "light" ? "light" : "dark";

  const fileContentQuery = useQuery({
    queryKey: ['devenv', devenvUuid, 'files', filename, 'content'],
    queryFn: () => axios.get<string>(
      (process.env.NEXT_PUBLIC_API ?? "") + "/api/devenv/" + devenvUuid + "/files/" + filename,
      {
        withCredentials: true
      }
    ).then((data) => {
      return data.data
    }),
    staleTime: Infinity,
    enabled: Boolean(filename)
  })

  const fileContentMutation = useMutation({
    mutationFn: (value: string) => axios.post(
      (process.env.NEXT_PUBLIC_API ?? "") + "/api/devenv/" + devenvUuid + "/files/" + filename,
      value,
      {
        withCredentials: true
      }
    ),
    onSuccess: (_, value) => {
      queryClient.setQueryData<string>(
        ['devenv', devenvUuid, 'files', filename, 'content'],
        () => {
          return value
        }
      )
    },
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
