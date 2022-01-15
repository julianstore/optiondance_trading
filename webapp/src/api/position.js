import $axios from "@/api/index";

export function positionDetail(positionId) {
    return  $axios.get('/v2/position/'+positionId,{
    });
}

export function physicalExercise(data) {
   return   $axios.post("/v1/exercise-request", data);
}

export function cashExercise(data){
    return  $axios.post("/v1/exercise-cash", data);
}