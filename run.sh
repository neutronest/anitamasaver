go get github.com/anaskhan96/soup
go get github.com/bradfitz/slice

DIRECTORY=./output

if [ -d "$DIRECTORY" ]; 
then
    echo "$DIRECTORY is already existed" # Control will enter here if $DIRECTORY exists.
else
    mkdir $DIRECTORY
fi

echo "Building metadata.."
cd ./metadata;go build
echo "Done..Building articledata..."
cd ../articledata;go build
echo "Done.."

echo "begin to saving article metadata..."
cd ../metadata; ./metadata
echo "Done.. begin to saving articles.."
echo "This step will cost few hours.. PLEASE STAND BY..."
cd ../articledata; ./articledata
echo "Done.."

