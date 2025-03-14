declare module '@uiw/react-codemirror' {
    import * as React from 'react';

    export interface CodeMirrorProps {
        value?: string;
        onChange?: (value: string) => void;
        height?: string;
        extensions?: any[];
        theme?: 'light' | 'dark' | string;
        readOnly?: boolean;
        placeholder?: string;
        // Добавьте другие пропсы, которые вы используете
    }

    const CodeMirror: React.FC<CodeMirrorProps>;
    export default CodeMirror;
}