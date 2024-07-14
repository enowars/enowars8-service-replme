import { Button } from "./ui/button"
import { Skeleton } from "./ui/skeleton"
import { Cross2Icon } from "@radix-ui/react-icons"
import { z } from "zod"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Form, FormControl, FormField, FormItem } from "./ui/form"
import { Input } from "./ui/input"
import { useDevenvFilesQuery } from "@/hooks/use-devenv-files-query"
import { useDevenvCreateFileMutation } from "@/hooks/use-devenv-create-file-mutation"
import { useDevenvDeleteFileMutation } from "@/hooks/use-devenv-delete-file-mutation"

type FileTreeProps = {
  className?: string
  devenvUuid: string
  selectedFile?: string
  setSelectedFile: (file: string | undefined) => void
}

const CreateFileFormSchema = z.object({
  name: z.string().min(1).max(30),
});

type CreateFileForm = z.infer<typeof CreateFileFormSchema>;

const FileTree: React.FC<FileTreeProps> = (props) => {
  const { className, devenvUuid, selectedFile, setSelectedFile } = props;

  const form = useForm<CreateFileForm>({
    resolver: zodResolver(CreateFileFormSchema),
    defaultValues: {
      name: "",
    }
  })

  const filesQuery = useDevenvFilesQuery({
    uuid: devenvUuid,
    callback: (files) => {
      if (!selectedFile && files.length > 0)
        setSelectedFile(files[0])
    }
  })

  const createFileMutation = useDevenvCreateFileMutation({
    uuid: devenvUuid,
    onSuccess: (file) => {
      if (!selectedFile)
        setSelectedFile(file.name)
    }
  })

  const deleteFileMutation = useDevenvDeleteFileMutation({
    uuid: devenvUuid,
    onSuccess: (filename, files) => {
      if (selectedFile === filename)
        setSelectedFile(files.length > 0 ? files[0] : undefined)
    }
  });

  const onCreateFileSubmit = (file: CreateFileForm) => {
    createFileMutation.mutate(file)
  }

  return (
    <div className={className}>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onCreateFileSubmit)} className="space-y-5">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input placeholder="filename" {...field} />
                </FormControl>
              </FormItem>
            )}
          />
          <Button type="submit" className="w-full mb-5">
            Add file
          </Button>
        </form>
      </Form>
      {filesQuery.isLoading && (
        <>
          <Skeleton className="w-full h-[20px] rounded-full" />
          <Skeleton className="w-full h-[20px] rounded-full" />
        </>
      )}
      {filesQuery.isSuccess &&
        <div key={"files-container"} className="flex flex-col px-4 w-full h-full overflow-y-scroll">
          {filesQuery.data.map((filename) => (
            <div key={"_" + filename} className={"flex flex-row justify-between rounded-full items-center pl-3 cursor-pointer " + (selectedFile === filename ? "bg-black text-white dark:bg-white dark:text-black " : "")} onClick={() => setSelectedFile(filename)}>
              <div>{filename}</div>
              <Button variant="ghost" className="rounded-full" onClick={(event) => {
                event.stopPropagation();
                deleteFileMutation.mutate(filename)
              }}>
                <Cross2Icon className="h-4 w-4" />
              </Button>
            </div>
          ))}
        </div>
      }
      {filesQuery.isError &&
        <div>Uh oh, something went wrong</div>
      }
    </div>
  )
}

export default FileTree;
