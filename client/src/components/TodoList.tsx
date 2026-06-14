// STARTER CODE:

import { Flex, Spinner, Stack, Text } from "@chakra-ui/react";
import TodoItem from "./TodoItem";
import { useQuery } from "@tanstack/react-query";
import { BASE_URL } from "../App";

export type Todo = {
    _id: number;
    body: string;
    completed: boolean;
}

const TodoList = () => {
	//const [isLoading, setIsLoading] = useState(false);
    const {data:todos, isLoading} =useQuery<Todo[]>({
        queryKey:["todos"],
        queryFn: async() => {
            try {
                const res = await fetch(BASE_URL + "/todos")
                const dataBack = await res.json()

                if (!res.ok) {
                    throw new Error(dataBack.error || "Something went wrong")
                }
                return dataBack || []
            } catch (error) {
                console.log(error)
            }
        }
    })
	// const todos = [
	// 	{
	// 		_id: 1,
	// 		body: "Buy groceries",
	// 		completed: true,
	// 	},
	// 	{
	// 		_id: 2,
	// 		body: "Walk the dog",
	// 		completed: false,
	// 	},
	// 	{
	// 		_id: 3,
	// 		body: "Do laundry",
	// 		completed: false,
	// 	},
	// 	{
	// 		_id: 4,
	// 		body: "Cook dinner",
	// 		completed: true,
	// 	},
	// ];
	return (
		<>
			<Text bgGradient='linear(to-l, #0780b8, #090cc2)' bgClip='text' fontSize={"4xl"} textTransform={"uppercase"} fontWeight={"bold"} textAlign={"center"} my={2}>
				Today's Tasks
			</Text>
			{isLoading && (
				<Flex justifyContent={"center"} my={4}>
					<Spinner size={"xl"} />
				</Flex>
			)}
            {/* //Todos?.length will only check length of todos if todos is not null/undefined (todos === null ? undefined : todos.length) */}
			{!isLoading && todos?.length === 0 && (
				<Stack alignItems={"center"} gap='3'>
					<Text fontSize={"xl"} textAlign={"center"} color={"gray.500"}>
						All tasks completed! 🤞
					</Text>
					<img src='/go.png' alt='Go logo' width={70} height={70} />
				</Stack>
			)}
			<Stack gap={3}>
				{todos?.map((todo) => (
					<TodoItem key={todo._id} todo={todo} />
				))}
			</Stack>
		</>
	);
};
export default TodoList;