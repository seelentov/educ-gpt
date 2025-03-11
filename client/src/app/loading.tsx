import { Loading } from '../components/ui/loading';

export default async function LoadingPage() {

    return (
        <div className="col-12 d-flex justify-content-center align-items-center gap-5" style={{ height: 200 }}>
            <Loading />
        </div>
    );
}