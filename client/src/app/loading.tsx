import { Loading } from './loading/loading';

export default async function LoadingPage() {

    return (
        <div className="col-12 d-flex justify-content-center align-items-center" style={{ height: 200 }}>
            <Loading />
        </div>
    );
}