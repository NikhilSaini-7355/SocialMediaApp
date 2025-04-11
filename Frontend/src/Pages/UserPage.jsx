import UserHeader from "../Components/UserHeader";
import UserPost from "../Components/UserPost";

const UserPage = ()=>{
    return <>
    <UserHeader />
    <UserPost postTitle={"title1"} likes={10} replies={100} postImg={'./post1.png'}/>
    <UserPost postTitle={"title2"} likes={20} replies={200} postImg={'./post2.png'}/>
    <UserPost postTitle={"title3"} likes={30} replies={300} postImg={'./post3.png'}/>
    <UserPost postTitle={"title4"} likes={40} replies={400} postImg={''} />
    </>
}

export default UserPage;