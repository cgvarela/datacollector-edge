package textrecord

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/streamsets/datacollector-edge/api"
	"github.com/streamsets/datacollector-edge/api/fieldtype"
	"github.com/streamsets/datacollector-edge/container/recordio"
	"io"
)

type TextWriterFactoryImpl struct {
	// TODO: Add needed configs
}

func (t *TextWriterFactoryImpl) CreateWriter(
	context api.StageContext,
	writer io.Writer,
) (recordio.RecordWriter, error) {
	var recordWriter recordio.RecordWriter
	recordWriter = newRecordWriter(context, writer)
	return recordWriter, nil
}

type TextWriterImpl struct {
	context api.StageContext
	writer  *bufio.Writer
}

func (textWriter *TextWriterImpl) WriteRecord(r api.Record) error {
	textFieldValue, err := textWriter.getTextFieldPathValue(r.Get())
	if err != nil {
		return err
	}
	fmt.Fprintln(textWriter.writer, textFieldValue)
	return nil
}

func (textWriter *TextWriterImpl) getTextFieldPathValue(field api.Field) (string, error) {
	var textFieldValue string
	if field.Value == nil {
		return textFieldValue, nil
	}
	var err error = nil
	switch field.Type {
	case fieldtype.MAP:
		fieldValue := field.Value.(map[string]api.Field)
		textField := fieldValue["text"]
		if textField.Type != fieldtype.STRING {
			err = errors.New("Invalid Field Type for Text Field path - " + textField.Type)
			return textFieldValue, err
		}
		textFieldValue = textField.Value.(string)
		return textFieldValue, err
	default:
		err = errors.New("Unsupported Field Type")
	}
	return textFieldValue, err
}

func (textWriter *TextWriterImpl) Flush() error {
	return recordio.Flush(textWriter.writer)
}

func (textWriter *TextWriterImpl) Close() error {
	return recordio.Close(textWriter.writer)
}

func newRecordWriter(context api.StageContext, writer io.Writer) *TextWriterImpl {
	return &TextWriterImpl{
		context: context,
		writer:  bufio.NewWriter(writer),
	}
}
