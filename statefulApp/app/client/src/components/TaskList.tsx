import {useEffect, useRef, useState} from "react";
import {fetchAPI} from "@/components/network/network.ts";
import {useNavigate} from "react-router";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card.tsx";
import {Button} from "@/components/ui/button.tsx";
import {CheckIcon} from "lucide-react";
import {Input} from "@/components/ui/input.tsx";
import {DatePicker} from "@/components/DatePicker.tsx";
import {format} from "date-fns";

export const TaskList = () => {
    const navigate = useNavigate();
    const [tasks, setTasks] = useState([]);
    const form = useRef(null);

    useEffect(() => {
        getTasks()
    }, []);

    const getTasks = async () => {
        const response = await fetchAPI("/api/tasks");
        if (!response.ok) {
            const data = await response.json()
            navigate("/", {
                state: {
                    error: data.error,
                    message: data.message,
                    status: response.status,
                }
            });
        }
        const data = await response.json();
        setTasks(data.tasks);
    }

    const handleCreateTask = async () => {
        if (!form.current) return;
        const currentForm = form.current as HTMLFormElement;

        const formData = new FormData(currentForm);
        const taskData = {
            name: formData.get("task_name"),
            description: formData.get("description"),
            due_date: formData.get("date"),
        };
        const previousTasks = tasks;

        // @ts-ignore
        setTasks([...tasks, {...taskData, id: -1}]);

        const response = await fetchAPI("/api/tasks", {
            method: "POST",
            body: JSON.stringify(taskData)
        });

        if (response.ok) {
            currentForm.reset();
            const newTask = await response.json();

            // @ts-ignore
            setTasks([...previousTasks, newTask.task]);
        } else {
            setTasks(previousTasks)
        }
    };

    const handleRemoveTask = async (taskId: number) => {
        const response = await fetchAPI(`/api/tasks/${taskId}`, {
            method: "DELETE"
        });

        if (response.ok) {
            setTasks(tasks.filter((task: any) => task.id !== taskId))
        }
    };

    return (
        <div className={"flex flex-col justify-center items-center h-screen gap-4"}>
            <Card className={"min-w-[450px] w-3/5"}>
                <CardHeader>
                    <CardTitle>Your current tasks</CardTitle>
                </CardHeader>
                <CardContent className={"flex flex-col gap-2"}>
                    {tasks.map((task: any, index) => {
                        return (
                            <div key={task.id}
                                 className={`group flex flex-row justify-between items-center ${index < (tasks.length - 1) ? "border-b border-gray-200 py-2" : ""}`}>
                                <div>
                                    <h2 className={"text-xl font-bold"}>{task.name}</h2>
                                    <p>{task.description}</p>
                                </div>
                                <div className={"flex flex-row gap-2 items-center"}>
                                    <p className={"text-slate-500"}>{format(task.due_date, "PPP")}</p>
                                    <div className={"min-w-[80px] flex justify-end"}>
                                        <Button
                                            className={"group-hover:block hidden"}
                                            variant={"outline"}
                                            onClick={() => handleRemoveTask(task.id)}
                                        >
                                            <CheckIcon/>
                                        </Button>
                                    </div>
                                </div>
                            </div>
                        )
                    })}
                </CardContent>
            </Card>

            <Card className={"min-w-[450px] w-3/5"}>
                <CardHeader>
                    <CardTitle>What should you do next ?</CardTitle>
                </CardHeader>
                <form onSubmit={(e) => e.preventDefault()} ref={form}>
                    <CardContent className={"flex flex-col gap-2"}>
                        <div className={"flex flex-row gap-2"}>
                            <Input name={"task_name"} placeholder={"Task name"}/>
                            <DatePicker name={"date"}/>
                        </div>
                        <Input name={"description"} placeholder={"Task description"}/>
                        <Button onClick={handleCreateTask}>Create task</Button>
                    </CardContent>
                </form>
            </Card>
        </div>
    )
}