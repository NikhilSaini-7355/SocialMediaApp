import { Flex } from "@chakra-ui/react";


// some issues may be there with this file 

const Actions = ({liked,setLiked})=>{
    return (
        <Flex gap={3} my={2} onClick={(e)=> e.preventDefault()}>
            <svg
             aria-label="like"
             color={liked?"rgb(237,73,86)":""}
             fill = {liked? "rgb(237,73,86)":"transparent"}
             height={'19'}
             role="img"
             viewBox="0 0 24 22"
             width={'20'}
             onClick={()=>setLiked(!liked)}
             >
                <path
                  d='M1 7.66c0 4.575 3.899 9.086 9.987 12.934.338.203.74.406 1.013.406.283 0 .686-.203 1.013-.406C19.1 16.746 23 12.234 23 7.66 23 3.736 20.245 1 16.672 1 14.603 1 12.98 1.94 12 3.352 11.042 1.952 9.408 1 7.328 1 3.766 1 1 3.736 1 7.66Z'
                  stroke='currentColor'
                  strokeWidth='2'
                >    
                </path>
             </svg>

             <svg
             aria-label="comment"
             color=""
             fill=""
             height={'19'}
             role="img"
             viewBox="0 0 24 24"
             width={'20'}
             >
                <title>comment</title>
                <path
                   d='M20.656 17.008a9.993 9.993 0 1 0-3.59 3.615L22 22Z'
                   fill='none'
                   stroke='currentColor'
                   strokeLinejoin='round'
                   strokeWidth='2'
                >  
                </path>
             </svg>

             <svg
             aria-label="repost"
             height={'19'}
             role="img"
             viewBox="0 0 24 24"
             width={'20'}
             >
                <title>repost</title>
                <path
                   d='M20.656 17.008a9.993 9.993 0 1 0-3.59 3.615L22 22Z'
                   fill='none'
                   stroke='currentColor'
                   strokeLinejoin='round'
                   strokeWidth='2'
                >
                    
                </path>
             </svg>

             <svg
             aria-label="share"
             height={'19'}
             role="img"
             viewBox="0 0 24 24"
             width={'20'}
             >
                <title>share</title>
                <path
                   d='M20.656 17.008a9.993 9.993 0 1 0-3.59 3.615L22 22Z'
                   fill='none'
                   stroke='currentColor'
                   strokeLinejoin='round'
                   strokeWidth='2'
                >
                    
                </path>
             </svg>
        </Flex>
    )
}

export default Actions;