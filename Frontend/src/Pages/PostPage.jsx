import { Avatar, Box, Button, Divider, Flex, Image, Text } from "@chakra-ui/react";
import { BsThreeDots } from "react-icons/bs";
import Actions from "../Components/Actions";
import { useState } from "react";
import Comment from "../Components/Comment";

const PostPage = ()=>{
    const [liked,setLiked] = useState(false);
    return (
    <>
      <Flex>
        <Flex w={'full'} alignItems={'center'} gap={3}>
            <Avatar src="/zuck-avatar.png" size={'md'} name="Mark Zuckerberg" />
            <Flex>
                <Text fontSize={'sm'} fontWeight={'bold'}>
                    markzuckerberg
                </Text>
                <Image src='/verified.png' w={4} h={4} ml={4} />
            </Flex>
        </Flex>
        <Flex gap={4} alignItems={'center'}>
            <Text fontSize={'sm'} color={'gray.light'}>
               1d
            </Text>
            <BsThreeDots />
        </Flex>
      </Flex>
      <Text my={3}>Post1</Text>
      <Box borderRadius={6} overflow={'hidden'} border={'1px solid'} borderColor={'gray.light'}>
        <Image src='/post1.png' w={'full'} />
      </Box>
      <Flex gap={3} my={3}>
        <Actions liked={liked} setLiked={setLiked} />
      </Flex>

      <Flex gap={2} alignItems={'center'}>
        <Text color={'gray.light'} fontSize={'sm'}>238 replies</Text>
        <Box w={0.5} h={0.5} borderRadius={'full'} bg={'gray.light'}></Box>
        <Text color={'gray.light'} fontSize={'sm'}>
            {200 + (liked?1:0)} likes
        </Text>
      </Flex>
      <Divider my={4}/>

      <Flex justifyContent={'space-between'}>
        <Flex gap={2} alignItems={'center'}>
            <Text color={'gray.light'}>Get the app</Text>
        </Flex>
        <Button>Get</Button>
      </Flex>
      <Divider my={4} />

      <Comment comment={"hey this looks good"} userAvatar={'https://bit.ly/dan-abramov'} username={"johndoe"} likes={200} createdAt={'2d'} />
      <Comment comment={"hey bro this looks good"} userAvatar={'https://bit.ly/sage-adebayo'} username={"johndoe"} likes={200} createdAt={'2d'} />
      <Comment comment={"kal jaldi kyu soya?"} userAvatar={'https://bit.ly/prosper-baba'} username={"johndoe"} likes={200} createdAt={'2d'} />
      
    </>
    )
}

export default PostPage;