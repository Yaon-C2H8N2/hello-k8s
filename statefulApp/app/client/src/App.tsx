import {useLocation, useNavigate} from "react-router";
import {useEffect, useRef, useState} from "react";
import {Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle} from "@/components/ui/card.tsx";
import {Button} from "@/components/ui/button.tsx";
import {Input} from "@/components/ui/input.tsx";
import {fetchAPI} from "@/components/network/network.ts";
import Cookie from "js-cookie";

function App() {
    const {state} = useLocation();
    const navigate = useNavigate();
    const [signUpMode, setSignUpMode] = useState(false);
    const form = useRef(null);

    useEffect(() => {
        if (state) {
            console.log(state);
        }
    }, [state]);

    const handleSignIn = async () => {
        if(!form.current) {
            return;
        }
        const formData = new FormData((form.current as HTMLFormElement));

        const response = await fetchAPI("/api/authenticate", {
            method: "POST",
            body: JSON.stringify({
                username: formData.get("username"),
                password: formData.get("password"),
            }),
        })

        if(response.status === 200){
            const data = await response.json();
            Cookie.set("token", data.token);
            navigate("/tasks");
        }
    }

    const handleSignUp = async () => {
        if (signUpMode) {
            await fetchAPI("/api/register", {
                method: "POST",
                body: JSON.stringify({
                    username: "test",
                    password: "test",
                }),
            })
        } else {
            setSignUpMode(true);
        }
    }

    return (
        <div className={"flex justify-center items-center h-screen"}>
            <Card className={"w-[450px]"}>
                <CardHeader>
                    <CardTitle>Hello-k8s</CardTitle>
                    <CardDescription>Login to your account</CardDescription>
                </CardHeader>
                <form ref={form}>
                    <CardContent className={"flex flex-col gap-2"}>
                        <Input placeholder={"Username"} name={"username"}/>
                        <Input placeholder={"Password"} type={"password"} name={"password"}/>
                        {signUpMode && (
                            <>
                                <Input placeholder={"Confirm Password"} type={"password"}/>
                                <div className={"flex flex-row gap-2"}>
                                    <Input placeholder={"First Name"}/>
                                    <Input placeholder={"Last Name"}/>
                                </div>
                                <Input placeholder={"Age"} type={"number"}/>
                            </>
                        )}
                    </CardContent>
                </form>
                <CardFooter className={"flex justify-between"}>
                    <Button variant={"outline"} onClick={handleSignUp}>Sign Up</Button>
                    {!signUpMode && (<Button onClick={handleSignIn}>Sign In</Button>)}
                </CardFooter>
            </Card>
        </div>
    )
}

export default App
