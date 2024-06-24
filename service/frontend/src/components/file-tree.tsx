import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import axios from "axios"
import { Button } from "./ui/button"
import { Skeleton } from "./ui/skeleton"
import { Cross2Icon } from "@radix-ui/react-icons"
import { z } from "zod"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { Form, FormControl, FormField, FormItem } from "./ui/form"
import { Input } from "./ui/input"

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

  const queryClient = useQueryClient();

  const form = useForm<CreateFileForm>({
    resolver: zodResolver(CreateFileFormSchema),
    defaultValues: {
      name: "",
    }
  })

  const filesQuery = useQuery({
    queryKey: ['devenv', devenvUuid, 'files'],
    queryFn: () => axios.get<string[]>(
      process.env.NEXT_PUBLIC_API + "/api/devenv/" + devenvUuid + "/files",
      {
        withCredentials: true
      }
    ).then((data) => {
      const files = data.data.sort();
      if (!selectedFile && files.length > 0)
        setSelectedFile(files[0])
      return files
    }),
  })

  const createFileMutation = useMutation({
    mutationFn: (file: CreateFileForm) => axios.post(
      process.env.NEXT_PUBLIC_API + '/api/devenv/' + devenvUuid + "/files",
      file,
      {
        withCredentials: true
      }
    ),
    onSuccess: (_, file) => {
      queryClient.setQueryData<string[]>(
        ['devenv', devenvUuid, 'files'],
        (oldData) => {
          let _data: string[] = oldData ?? []
          if (_data.includes(file.name))
            return _data
          if (!selectedFile)
            setSelectedFile(file.name)
          return [..._data, file.name].sort()
        }
      )
    },
  })

  const deleteFileMutation = useMutation({
    mutationFn: (filename: string) => axios.delete(
      process.env.NEXT_PUBLIC_API + '/api/devenv/' + devenvUuid + "/files/" + encodeURI(filename),
      {
        withCredentials: true
      }
    ),
    onSuccess: (_, filename) => {
      queryClient.setQueryData<string[]>(
        ['devenv', devenvUuid, 'files'],
        (oldData) => {
          let _data: string[] = oldData ?? []
          _data = _data.filter((name) => name !== filename)
          if (selectedFile === filename)
            setSelectedFile(_data.length > 0 ? _data[0] : undefined)
          return _data
        }
      )
    }
  })

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
