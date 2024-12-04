from openai import OpenAI
import json
from functions import modify_config, restart_service, apply_manifest, tools

client = OpenAI(
    api_key="xxxxxxxxxxx",
    base_url="https://api.apiyi.com/v1",
)

def run_conversation():
    # 输入操作指令:
    command = input("请输入您想操作的内容: ")

    # 定义 prompt
    messages = [
        {
            "role": "system",
            "content": "你是一个服务运维执行助手，你可以帮助用户来执行不同的运维操作，你可以调用多个函数来帮助用户完成任务",
        },
        {
            "role": "user",
            "content": command,
        },
    ]

    # 准备调用模型，并且指定预定义的function
    response = client.chat.completions.create(
        model="gpt-4o",
        messages=messages,
        tools=tools,
        tool_choice="auto",
    )
    response_message = response.choices[0].message
    tool_calls = response_message.tool_calls
    print("\nChatGPT want to call function: ", tool_calls)

    # 根据大模型返回的function，已经对应的参数，调用对应的函数
    if tool_calls is None:
        print("No function called")
    if tool_calls:
        # 定义可被映射的函数
        available_functions = {
            "modify_config": modify_config,
            "restart_service": restart_service,
            "apply_manifest": apply_manifest,
        }
        # 将assistant返回的结果补充到messages，以便下次对话时提交给大模型
        messages.append(response_message)
        for tool_call in tool_calls:
            function_name = tool_call.function.name
            function_to_call = available_functions[function_name]
            function_args = json.loads(tool_call.function.arguments)
            print(f"Function name: {function_name}, arguments: {function_args}")
            function_result = function_to_call(**function_args)
            print(f"Function exec result: {function_result}")

            # 把执行过程及反馈结果补充到messages，重新提交给大模型，进行下一步的收尾对话
            messages.append(
                {
                    "role": "tool",
                    "tool_call_id": tool_call.id,
                    "name": function_name,
                    "content": function_result,
                }
            )

        response = client.chat.completions.create(
            model="gpt-4o",
            messages=messages,
        )

        return response.choices[0].message.content

if __name__ == '__main__':
    print('LLM Res:', run_conversation())

