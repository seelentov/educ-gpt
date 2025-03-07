import Link from "next/link";

export default function Home() {
    return (
        <div>
            <div className="full-height">
                <div className="jumbotron">
                    <div className="container">
                        <h1 className="display-3">AI-Репетитор у вас под рукой</h1>
                        <p>Получи знания и практику, адаптированные специально для тебя! Выбирай тему, изучай детали с помощью AI-помощника и выполняй задачи для закрепления материала. Индивидуальный подход к обучению!</p>
                        <p><Link href="/login" className="btn btn-primary btn-lg" role="button">Войти</Link></p>
                    </div>
                </div>
            </div>
        </div>
    );
}
